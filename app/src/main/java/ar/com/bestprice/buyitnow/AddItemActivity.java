package ar.com.bestprice.buyitnow;

import android.content.Intent;
import android.graphics.Bitmap;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.view.View;
import android.widget.ArrayAdapter;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.Spinner;
import android.widget.TextView;

import com.google.android.gms.common.api.CommonStatusCodes;
import com.google.android.gms.vision.barcode.Barcode;


import java.io.UnsupportedEncodingException;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.ArrayList;
import java.util.Arrays;

import ar.com.bestprice.buyitnow.barcodereader.BarcodeCaptureActivity;
import ar.com.bestprice.buyitnow.dto.Item;

/**
 * Created by ivan on 01/04/16.
 */
public class AddItemActivity extends AppCompatActivity implements View.OnClickListener{

    private static final int RC_BARCODE_CAPTURE = 5;
    private final static int CAMERA_REQUEST = 10;

    @Override
    protected void onCreate(Bundle savedInstanceState) {

        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_add_item);
        Spinner spinner = (Spinner) findViewById(R.id.add_item_spinner);

        ArrayList<Category> arraySpinner = new ArrayList<>();

        arraySpinner.addAll(Arrays.asList(Category.values()));

        ArrayAdapter<Category> adapter = new ArrayAdapter<>(this,android.R.layout.simple_spinner_item, arraySpinner);

        spinner.setAdapter(adapter);

    }

    @Override
    public void onClick(View v) {

        switch(v.getId()) {

            case R.id.imageButton_barcode: //barcode button was pushed, so we need to capture the product barcode.

                Intent intent = new Intent(this.getApplicationContext(), BarcodeCaptureActivity.class);

                intent.putExtra(BarcodeCaptureActivity.AutoFocus, true);
                intent.putExtra(BarcodeCaptureActivity.UseFlash, false);

                startActivityForResult(intent, RC_BARCODE_CAPTURE);

                break;

            case R.id.save_purchase_item: //save item button was pushed

                EditText id = (EditText)findViewById(R.id.add_item_prod_id);
                EditText description = (EditText)findViewById(R.id.description_text);
                EditText price = (EditText)findViewById(R.id.price_text);

                Spinner spinner = (Spinner)findViewById(R.id.add_item_spinner);

                TextView textView = (TextView)spinner.getSelectedView();
                String category = textView.getText().toString();


                String itemId = id.getText().toString();

                if(itemId.isEmpty()){

                    MessageDigest crypt = null;
                    try {
                        crypt = MessageDigest.getInstance("SHA-1");
                        crypt.reset();
                        crypt.update(description.getText().toString().getBytes("UTF-8"));
                    } catch (NoSuchAlgorithmException e) {
                        e.printStackTrace();
                    } catch (UnsupportedEncodingException e) {
                        e.printStackTrace();
                    }

                    itemId = Context.byteToHex(crypt.digest());
                }

                Item item = new Item();
                item.setId(itemId);
                item.setDescription(description.getText().toString());
                item.setPrice(Float.valueOf(price.getText().toString()));
                item.setCategory(category);

                startActivity(item);
                finish();
                break;

            case R.id.take_picture:
                Intent cameraIntent = new Intent(android.provider.MediaStore.ACTION_IMAGE_CAPTURE);
                startActivityForResult(cameraIntent, CAMERA_REQUEST);
                break;
        }

    }

    private void startActivity(Item item){

        //If add item activity was called by AddNewPurchase activity, so the item has to be returned
        //as data
       if (getIntent().getIntExtra(Constants.CALLING_ACTIVITY, 0) == Constants.NEW_PURCHASE){

            Intent data = new Intent();
            data.putExtra(Constants.ITEM, item);
            setResult(CommonStatusCodes.SUCCESS, data);

       //If Add item activity was called by MainActivity, so AddNewPurchase activity is started
       //and the item is passes as parameter
       } else {

            Intent intent = new Intent(this.getApplicationContext(), AddPurchaseActivity.class);
            intent.putExtra(Constants.ITEM, item);
            startActivity(intent);
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {

        if(requestCode == RC_BARCODE_CAPTURE && resultCode == CommonStatusCodes.SUCCESS && data != null){

            Barcode barcode = data.getParcelableExtra(BarcodeCaptureActivity.BarcodeObject);
            EditText id = (EditText)findViewById(R.id.add_item_prod_id);
            id.setText(barcode.displayValue);

        }else if (requestCode == CAMERA_REQUEST && resultCode == RESULT_OK) {

            Bitmap photo = (Bitmap) data.getExtras().get("data");
            ImageView imageView = (ImageView)findViewById(R.id.product_picture);
            imageView.setImageBitmap(photo);
        }
    }
}
