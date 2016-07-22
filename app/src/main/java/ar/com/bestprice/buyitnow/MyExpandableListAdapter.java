package ar.com.bestprice.buyitnow;

import android.app.Activity;
import android.graphics.Color;
import android.graphics.drawable.Drawable;
import android.view.ActionMode;
import android.view.LayoutInflater;
import android.view.Menu;
import android.view.MenuInflater;
import android.view.MenuItem;
import android.view.View;
import android.view.ViewGroup;
import android.widget.BaseExpandableListAdapter;
import android.widget.CheckedTextView;
import android.widget.ExpandableListView;
import android.widget.LinearLayout;
import android.widget.RelativeLayout;
import android.widget.TextView;

import java.util.ArrayList;
import java.util.Map;

import ar.com.bestprice.buyitnow.dto.Item;
import ar.com.bestprice.buyitnow.dto.Purchase;
import ar.com.bestprice.buyitnow.dto.Purchases;


public class MyExpandableListAdapter extends BaseExpandableListAdapter {

    private final Map<Integer, PurchasesGroup> groups;
    public LayoutInflater inflater;
    public Activity activity;
    private ExpandableListView myParentExpandableView;

    public MyExpandableListAdapter(Activity act, Map<Integer, PurchasesGroup> groups) {
        activity = act;
        this.groups = groups;
        inflater = act.getLayoutInflater();
    }

    @Override
    public Object getChild(int groupPosition, int childPosition) {

        return groups.get(groupPosition).getItemAt(childPosition);
    }

    @Override
    public long getChildId(int groupPosition, int childPosition) {
        return 0;
    }


    @Override
    public View getChildView(final int groupPosition, final int childPosition,
                             boolean isLastChild, View convertView, final ViewGroup parent) {

        final Item children = (Item) getChild(groupPosition, childPosition);

        if (convertView == null) {
            convertView = inflater.inflate(R.layout.listrow_details, null);
        }

        TextView text = (TextView) convertView.findViewById(R.id.listrow_item_description);
        text.setText(children.getDescription());

        text.setOnLongClickListener(new View.OnLongClickListener() {

            @Override
            public boolean onLongClick(final View view) {
                final Drawable color = ((RelativeLayout)view.getParent().getParent()).getBackground();
                ((RelativeLayout)((view.getParent()).getParent())).setBackgroundColor(Color.LTGRAY);
                myParentExpandableView.startActionMode(new ActionMode.Callback() {// Called when the action mode is created; startActionMode() was called

                    @Override
                    public boolean onCreateActionMode(ActionMode mode, Menu menu) {
                        // Inflate a menu resource providing context menu items
                        MenuInflater inflater = mode.getMenuInflater();
                        inflater.inflate(R.menu.delete_item_menu, menu);
                        return true;
                    }

                    @Override
                    public boolean onPrepareActionMode(ActionMode mode, Menu menu) {
                        return false; // Return false if nothing is done
                    }

                    @Override
                    public boolean onActionItemClicked(ActionMode mode, MenuItem item) {
                        switch (item.getItemId()) {
                            case R.id.delete:
                                ((RelativeLayout)((view.getParent()).getParent())).setBackground(color);

                                PurchasesGroup group = (PurchasesGroup) getGroup(groupPosition);
                                group.removeItemAt(childPosition);

                                Purchase purchase = group.getPurchase(children.getTime());


                                mode.finish(); // Action picked, so close the CAB
                                notifyDataSetChanged();
                                return true;

                            default:
                                return false;
                        }
                    }

                    // Called when the user exits the action mode
                    @Override
                    public void onDestroyActionMode(ActionMode mode) {
                        ((RelativeLayout)((view.getParent()).getParent())).setBackground(color);
                    }
                });

                return true;
            }
        });

        int icon = Category.MERCADERIA.getIcon();
        if (children.getCategory() != null) {
            icon = children.getCategory().getIcon();
        }

        text.setCompoundDrawablesWithIntrinsicBounds(R.drawable.mercaderia_32, 0, 0, 0);
        text = (TextView) convertView.findViewById(R.id.item_price);
        text.setText(String.format("$%.2f", children.getPrice()));

        return convertView;
    }

    @Override
    public int getChildrenCount(int groupPosition) {
        return groups.get(groupPosition).getChildren().size();

    }

    @Override
    public Object getGroup(int groupPosition) {
        return groups.get(groupPosition);
    }

    @Override
    public int getGroupCount() {
        return groups.size();
    }

    @Override
    public void onGroupCollapsed(int groupPosition) {
        super.onGroupCollapsed(groupPosition);
    }

    @Override
    public void onGroupExpanded(int groupPosition) {
        super.onGroupExpanded(groupPosition);
    }

    @Override
    public long getGroupId(int groupPosition) {
        return 0;
    }

    @Override
    public View getGroupView(int groupPosition, boolean isExpanded,
                             View convertView, ViewGroup parent) {

        if (convertView == null) {
            convertView = inflater.inflate(R.layout.listrow_group, null);
        }

        RelativeLayout relativeLayout = (RelativeLayout) ((LinearLayout) convertView).getChildAt(0);

        CheckedTextView checkedTextView = (CheckedTextView) relativeLayout.getChildAt(0);
        TextView amountPerMonth = (TextView) relativeLayout.getChildAt(1);
        TextView differencePerMonth = (TextView) relativeLayout.getChildAt(3);

        PurchasesGroup purchasesByMonth = (PurchasesGroup) getGroup(groupPosition);
        PurchasesGroup previousPurchasesByMonth = (PurchasesGroup) getGroup(((groupPosition > 0) ? (groupPosition - 1) : 0));

        if (purchasesByMonth.getPurchasesTotalPrice() > previousPurchasesByMonth.getPurchasesTotalPrice()) {
            differencePerMonth.setBackgroundColor(Color.rgb(139, 00, 00)); //RED
        } else {
            differencePerMonth.setBackgroundColor(Color.rgb(34, 139, 34)); //GREEN
        }

        checkedTextView.setText(purchasesByMonth.getString());


        float diff = purchasesByMonth.getPurchasesTotalPrice() - previousPurchasesByMonth.getPurchasesTotalPrice();

        diff = (diff * 100) / previousPurchasesByMonth.getPurchasesTotalPrice();

        differencePerMonth.setText(String.format("%+.2f%%", diff));
        amountPerMonth.setText(String.format("$%.2f", purchasesByMonth.getPurchasesTotalPrice()));


        return convertView;
    }

    @Override
    public boolean hasStableIds() {
        return false;
    }

    @Override
    public boolean isChildSelectable(int groupPosition, int childPosition) {
        return true;
    }

    public void setParent(ExpandableListView parent) {
        this.myParentExpandableView = parent;
    }

}